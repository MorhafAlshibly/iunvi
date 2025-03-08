import {
  BasePlugin,
  Uppy,
  UppyFile,
  UploadResult,
  PluginOpts,
} from "@uppy/core";
import { DataLakeFileClient } from "@azure/storage-file-datalake";
import { AbortError } from "@azure/abort-controller";
import { FileProgressStarted } from "@uppy/utils/lib/FileProgress";

interface DataLakePluginOptions extends PluginOpts {
  getSasUrl: (file: UppyFile<{}, {}>) => Promise<string>;
  chunkSize: number;
  maxRetries: number;
}

export default class DataLakePlugin extends BasePlugin<
  DataLakePluginOptions,
  {},
  {}
> {
  #abortControllers: Map<string, AbortController>;
  #uploadHandler: (fileIDs: string[]) => Promise<void>;
  #fileRemovedHandler: (file: UppyFile<{}, {}>) => void;

  constructor(uppy: Uppy, opts: DataLakePluginOptions) {
    const defaultOptions = {
      getSasUrl: async () => {
        throw new Error("getSasUrl not provided");
      },
      chunkSize: 10 * 1024 * 1024, // 10 MB by default
      maxRetries: 3,
    };

    super(uppy, { ...defaultOptions, ...opts });
    this.id = opts.id || "DataLakePlugin";
    this.type = "uploader";

    this.#abortControllers = new Map();
    this.#uploadHandler = this.#uploadFiles.bind(this);
    this.#fileRemovedHandler = this.#stopUpload.bind(this);
  }

  /**
   * Upload multiple files concurrently.
   */
  async #uploadFiles(fileIDs: string[]): Promise<void> {
    this.uppy.on("file-removed", this.#fileRemovedHandler);

    await Promise.all(
      fileIDs.map((fileID) => {
        const file = this.uppy.getFile(fileID);
        return this.#uploadFileParallel(file);
      }),
    );

    this.uppy.off("file-removed", this.#fileRemovedHandler);
  }

  async #uploadFileParallel(file: UppyFile<{}, {}>): Promise<void> {
    const abortController = new AbortController();
    this.#abortControllers.set(file.id, abortController);

    this.uppy.emit("upload-start", [file]);

    try {
      await this.#startUploadParallel(file, (progress: FileProgressStarted) => {
        this.uppy.emit("upload-progress", file, progress);
      });
      this.uppy.emit("upload-success", file, { status: 200 });
    } catch (error: any) {
      this.uppy.emit("upload-error", file, {
        name: error.name,
        message: error.message,
      });
    } finally {
      this.#finishUpload(file);
    }
  }

  /**
   * Parallelizes chunk uploads for a single file.
   */
  async #startUploadParallel(
    file: UppyFile<{}, {}>,
    onProgress: (progress: FileProgressStarted) => void,
  ): Promise<void> {
    const { getSasUrl, chunkSize, maxRetries } = this.opts;

    if (!file.data || typeof file.data.slice !== "function") {
      throw new Error("Invalid file data; expected a Blob or File.");
    }

    const sasUrl = await getSasUrl(file);
    const dataLakeClient = new DataLakeFileClient(sasUrl);
    const totalSize = file.data.size;
    const totalChunks = Math.ceil(totalSize / chunkSize);

    // Array to track progress for each chunk.
    const chunkProgress: number[] = new Array(totalChunks).fill(0);
    const updateOverallProgress = () => {
      const uploadedBytes = chunkProgress.reduce((a, b) => a + b, 0);
      onProgress({
        uploadStarted: file.progress.uploadStarted ?? 0,
        bytesUploaded: uploadedBytes,
        bytesTotal: totalSize,
      });
    };

    // Function to upload a single chunk.
    const uploadChunk = async (
      i: number,
    ): Promise<{ index: number; blockId: string }> => {
      const start = i * chunkSize;
      const end = Math.min(start + chunkSize, totalSize);
      const chunk = file.data.slice(start, end);
      // Generate a fixed-length block ID (e.g. "000001", "000002", …)
      const blockId = btoa(i.toString().padStart(6, "0"));

      let attempt = 0;
      // Callback for progress for this chunk.
      const onChunkProgress = (progress: { loadedBytes: number }) => {
        chunkProgress[i] = progress.loadedBytes;
        updateOverallProgress();
      };

      while (attempt < maxRetries) {
        try {
          await dataLakeClient.stageChunk(blockId, chunk, chunk.size, {
            abortSignal: this.#abortControllers.get(file.id)?.signal,
            onProgress: onChunkProgress,
          });
          return { index: i, blockId };
        } catch (error: any) {
          attempt++;
          if (attempt >= maxRetries) {
            throw new Error(
              `Failed to stage block ${blockId} after ${maxRetries} attempts: ${error.message}`,
            );
          }
        }
      }
      // Unreachable; needed to satisfy TS.
      throw new Error("Unexpected error during chunk upload.");
    };

    // Launch all chunk uploads in parallel.
    const uploadPromises: Promise<{ index: number; blockId: string }>[] = [];
    for (let i = 0; i < totalChunks; i++) {
      uploadPromises.push(uploadChunk(i));
    }

    // Wait for all chunks to upload.
    const results = await Promise.all(uploadPromises);
    // Ensure the block IDs are ordered correctly.
    const sortedResults = results.sort((a, b) => a.index - b.index);
    const blockIds = sortedResults.map((result) => result.blockId);

    // Commit all staged blocks to finalize the blob.
    await blockBlobClient.commitBlockList(blockIds, {
      abortSignal: this.#abortControllers.get(file.id)?.signal,
    });
  }

  #stopUpload(file: UppyFile<{}, {}>): void {
    this.#abortControllers.get(file.id)?.abort();
  }

  #finishUpload(file: UppyFile<{}, {}>): void {
    this.#abortControllers.delete(file.id);
  }

  install(): void {
    this.uppy.addUploader(this.#uploadHandler);
  }

  uninstall(): void {
    this.uppy.removeUploader(this.#uploadHandler);
  }
}
