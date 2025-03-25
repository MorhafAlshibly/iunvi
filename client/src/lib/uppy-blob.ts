import {
  BasePlugin,
  Uppy,
  UppyFile,
  UploadResult,
  PluginOpts,
} from "@uppy/core";
import { BlockBlobClient } from "@azure/storage-blob";
import { AbortError } from "@azure/abort-controller";
import { FileProgressStarted } from "@uppy/utils/lib/FileProgress";

interface BlobPluginOptions extends PluginOpts {
  getSasUrl: (file: UppyFile<{}, {}>) => Promise<string>;
}

export default class BlobPlugin extends BasePlugin<BlobPluginOptions, {}, {}> {
  #abortControllers: Map<string, AbortController>;
  #uploadHandler: (fileIDs: string[]) => Promise<void>;
  #fileRemovedHandler: (file: UppyFile<{}, {}>) => void;

  constructor(uppy: Uppy, opts: BlobPluginOptions) {
    const defaultOptions = {
      getSasUrl: async () => {
        throw new Error("getSasUrl not provided");
      },
    };

    super(uppy, { ...defaultOptions, ...opts });
    this.id = opts.id || "BlobPlugin";
    this.type = "uploader";

    this.#abortControllers = new Map();
    this.#uploadHandler = this.#uploadFiles.bind(this);
    this.#fileRemovedHandler = this.#stopUpload.bind(this);
  }

  /**
   * Uploads multiple files concurrently.
   */
  async #uploadFiles(fileIDs: string[]): Promise<void> {
    this.uppy.on("file-removed", this.#fileRemovedHandler);

    const files = fileIDs.map((fileID) => this.uppy.getFile(fileID));
    this.uppy.emit("upload-start", files);

    for (const file of files) {
      const abortController = new AbortController();
      this.#abortControllers.set(file.id, abortController);

      try {
        await this.#startUpload(file, (progress) => {
          this.uppy.emit("upload-progress", file, progress);
        });

        this.uppy.emit("upload-success", file, {
          status: 200,
        });
      } catch (error: any) {
        this.uppy.emit("upload-error", file, {
          name: error.name,
          message: error.message,
        });
      } finally {
        this.#finishUpload(file);
      }
    }

    this.uppy.off("file-removed", this.#fileRemovedHandler);
  }

  /**
   * Starts the upload of a single file.
   */
  async #startUpload(
    file: UppyFile<{}, {}>,
    onProgress: (progress: FileProgressStarted) => void,
  ): Promise<void> {
    // Validate that file.data is a Blob/File-like object.
    if (!file.data || typeof file.data.slice !== "function") {
      throw new Error("Invalid file data; expected a Blob or File.");
    }

    try {
      // Obtain a SAS URL (with write permissions) for this file.
      const sasUrl = await this.opts.getSasUrl(file);
      const blockBlobClient = new BlockBlobClient(sasUrl);
      await blockBlobClient.uploadData(file.data, {
        abortSignal: this.#abortControllers.get(file.id)?.signal,
        onProgress: (progress) => {
          onProgress({
            uploadStarted: file.progress.uploadStarted ?? 0,
            bytesUploaded: progress.loadedBytes,
            bytesTotal: file.data.size,
          });
        },
      });
    } catch (error: any) {
      if (error.name === "AbortError") return;
      if (error.code === "UnauthorizedBlobOverwrite") {
        throw new Error("File already exists.");
      }
      throw error;
    }
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
