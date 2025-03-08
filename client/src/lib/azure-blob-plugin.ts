import {
  BasePlugin,
  Uppy,
  UppyFile,
  UploadResult,
  PluginOpts,
  DefinePluginOpts,
  Body,
  Meta,
} from "@uppy/core";
import { BlockBlobClient } from "@azure/storage-blob";
import { AbortError } from "@azure/abort-controller";
import { FileProgressStarted } from "@uppy/utils/lib/FileProgress";

/**
 * Options for the AzureBlobPlugin plugin.
 */
interface AzureBlobPluginOptions extends PluginOpts {
  getSasUrl: (file: UppyFile<{}, {}>) => Promise<string>;
  chunkSize: number;
  maxRetries: number;
}

export default class AzureBlobPlugin extends BasePlugin<
  AzureBlobPluginOptions,
  {},
  {}
> {
  #abortControllers;
  #uploadHandler;
  #fileRemovedHandler;

  constructor(uppy: Uppy, opts: AzureBlobPluginOptions) {
    const defaultOptions = {
      getSasUrl: async () => {
        throw new Error("getSasUrl not provided");
      },
      chunkSize: 10 * 1024 * 1024, // 10 MB by default
      maxRetries: 3,
    };

    super(uppy, { ...defaultOptions, ...opts });
    this.id = opts.id || "AzureBlobPlugin";
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

    for (const fileID of fileIDs) {
      const file = this.uppy.getFile(fileID);
      const abortController = new AbortController();
      this.#abortControllers.set(file.id, abortController);

      this.uppy.emit("upload-start", [file]);

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
    const { getSasUrl, chunkSize, maxRetries } = this.opts;

    // Validate that file.data is a Blob/File-like object.
    if (!file.data || typeof file.data.slice !== "function") {
      throw new Error("Invalid file data; expected a Blob or File.");
    }

    try {
      // Obtain a SAS URL (with write permissions) for this file.
      const sasUrl = await getSasUrl(file);
      const blockBlobClient = new BlockBlobClient(sasUrl);
      const totalSize = file.data.size;
      const totalChunks = Math.ceil(totalSize / chunkSize);
      let uploadedBytes = 0;
      const blockIds: string[] = [];

      for (let i = 0; i < totalChunks; i++) {
        const start = i * chunkSize;
        const end = Math.min(start + chunkSize, totalSize);
        const chunk = file.data.slice(start, end);
        // Create a unique block ID using base64 encoding.
        const blockIdRaw = `block-${i}`;
        const blockId = btoa(blockIdRaw);

        // Stage (upload) the block with retry logic.
        let attempt = 0;
        while (attempt < maxRetries) {
          try {
            await blockBlobClient.stageBlock(blockId, chunk, chunk.size, {
              abortSignal: this.#abortControllers.get(file.id)?.signal,
              onProgress: (progress) => {
                onProgress({
                  uploadStarted: file.progress.uploadStarted ?? 0,
                  bytesUploaded: progress.loadedBytes + uploadedBytes,
                  bytesTotal: totalSize,
                });
              },
            });
            break;
          } catch (error: any) {
            attempt++;
            if (attempt >= maxRetries) {
              throw new Error(
                `Failed to stage block ${blockId} after ${maxRetries} attempts: ${error.message}`,
              );
            }
          }
        }
        blockIds.push(blockId);
        uploadedBytes += chunk.size;
      }

      // Commit all staged blocks to finalize the blob.
      await blockBlobClient.commitBlockList(blockIds);
    } catch (error: any) {
      if (!(error instanceof AbortError)) {
        throw error;
      }
    }
  }

  #stopUpload(file: UppyFile<{}, {}>) {
    this.#abortControllers.get(file.id)?.abort();
  }

  #finishUpload(file: UppyFile<{}, {}>) {
    this.#abortControllers.delete(file.id);
  }

  /**
   * Install the plugin by adding an uploader hook.
   */
  install(): void {
    // Uppy calls this uploader function with an array of file IDs.
    this.uppy.addUploader(this.#uploadHandler);
  }

  /**
   * Uninstall the plugin. (Cleanup logic can be added here if needed.)
   */
  uninstall(): void {
    this.uppy.removeUploader(this.#uploadHandler);
  }
}

// import { Uppy, Plugin, PluginOpts } from "@uppy/core";
// import {
//   BlobServiceClient,
//   BlockBlobParallelUploadOptions,
// } from "@azure/storage-blob";
// import { AbortError } from "@azure/abort-controller";

// interface AzureBlobOptions extends PluginOpts {
//   endpoint: string;
//   container: string;
//   sas: string;
// }

// // declare class AzureBlobPlugin extends Plugin<AzureBlobOptions> {}

// export default class AzureBlobPlugin extends Plugin<
//   AzureBlobOptions,
//   {},
//   {},
//   {}
// > {
//   #abortControllers;
//   #containerClient;
//   #uploadHandler;
//   #fileRemovedHandler;

//   constructor(uppy: Uppy, opts: AzureBlobOptions) {
//     super(uppy, opts);
//     this.id = opts.id || "AzureBlobPlugin";
//     this.type = "uploader";

//     this.#abortControllers = new Map();

//     const blobServiceClient = new BlobServiceClient(opts.endpoint + opts.sas);
//     this.#containerClient = blobServiceClient.getContainerClient(
//       opts.container,
//     );

//     this.#uploadHandler = this.uploadFiles.bind(this);
//     this.#fileRemovedHandler = this.#stopUpload.bind(this);
//   }

//   install() {
//     this.uppy.addUploader(this.#uploadHandler);
//   }

//   uninstall() {
//     this.uppy.removeUploader(this.#uploadHandler);
//   }

//   async uploadFiles(fileIDs: any) {
//     this.uppy.on("file-removed", this.#fileRemovedHandler);

//     for (const fileID of fileIDs) {
//       const file = this.uppy.getFile(fileID);

//       try {
//         this.uppy.emit("upload", fileID, [file]);

//         const upload = await this.#startUpload([file], (progress: any) => {
//           this.uppy.emit("upload-progress", file, progress);
//         });

//         this.uppy.emit("upload-success", file, upload);
//       } catch (error) {
//         this.uppy.emit("upload-error", file, error);
//       } finally {
//         this.#finishUpload(file);
//       }
//     }

//     this.uppy.off("file-removed", this.#fileRemovedHandler);
//   }

//   async #startUpload(
//     file: {
//       id: any;
//       name: string;
//       data: any;
//       blobOptions: BlockBlobParallelUploadOptions | undefined;
//     },
//     onProgress: (progress: any) => void,
//   ) {
//     const abortController = new AbortController();
//     this.#abortControllers.set(file.id, abortController);

//     const blockBlobClient = this.#containerClient.getBlockBlobClient(file.name);

//     try {
//       return await blockBlobClient.uploadData(file.data, {
//         ...file.blobOptions,
//         abortSignal: abortController.signal,
//         onProgress: (progress) =>
//           onProgress(this.#azureProgressToUppyProgress(progress, file)),
//       });
//     } catch (error) {
//       if (!(error instanceof AbortError)) {
//         throw error;
//       }
//     }
//   }

//   #azureProgressToUppyProgress(
//     azureProgress: { loadedBytes: any },
//     file: { size: any },
//   ) {
//     return {
//       uploader: this,
//       bytesUploaded: azureProgress.loadedBytes,
//       bytesTotal: file.size,
//     };
//   }

//   #stopUpload(file: { id: any }) {
//     this.#abortControllers.get(file.id)?.abort();
//   }

//   #finishUpload(file: { id: any }) {
//     this.#abortControllers.delete(file.id);
//   }
// }
