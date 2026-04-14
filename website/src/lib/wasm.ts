declare global {
  interface Window {
    Go: new () => { run: (instance: WebAssembly.Instance) => void; importObject: WebAssembly.Imports };
    shiftcomply: {
      jurisdictions: () => string;
      rules: (jurisdiction: string, staff?: string, unit?: string, scope?: string) => string;
      constraints: (jurisdiction: string, staff?: string, unit?: string, scope?: string) => string;
      compare: (left: string, right: string, staff?: string) => string;
      validate: (scheduleJSON: string) => string;
      export: (jurisdiction: string) => string;
    };
  }
}

let loaded = false;
let loading: Promise<void> | null = null;

export function loadWasm(): Promise<void> {
  if (loaded) return Promise.resolve();
  if (loading) return loading;

  loading = new Promise<void>(async (resolve, reject) => {
    try {
      const basePath = process.env.NODE_ENV === "production" ? "/shift-comply" : "";

      // Load wasm_exec.js
      await new Promise<void>((res, rej) => {
        const script = document.createElement("script");
        script.src = `${basePath}/wasm_exec.js`;
        script.onload = () => res();
        script.onerror = () => rej(new Error("Failed to load wasm_exec.js"));
        document.head.appendChild(script);
      });

      const go = new window.Go();
      const result = await WebAssembly.instantiateStreaming(
        fetch(`${basePath}/shiftcomply.wasm`),
        go.importObject
      );
      go.run(result.instance);
      loaded = true;
      resolve();
    } catch (err) {
      reject(err);
    }
  });

  return loading;
}

export function isLoaded(): boolean {
  return loaded;
}
