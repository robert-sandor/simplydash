import type {Config} from "./config";
import type {Category} from "./category";

export class ApiClient {
    host: string

    constructor(host: string) {
        this.host = host;
    }

    async getConfig(): Promise<Config> {
        const response = await fetch(`${this.host}/api/config`);
        return await response.json();
    }

    async getCategories(): Promise<Category[]> {
        const response = await fetch(`${this.host}/api/categories`);
        return await response.json();
    }

    async healthcheck(url: string): Promise<number> {
        const response = await fetch(`${this.host}/api/url/health?url=${encodeURIComponent(url)}`);
        return response.status;
    }

    async websocket(
        onUpdateConfig: () => void,
        onUpdateCategories: () => void
    ) {
        const ws = new WebSocket(`ws://localhost:8080/api/ws`)

        ws.onmessage = (e: MessageEvent<string>) => {
            if (e.data === "update-config") {
                onUpdateConfig();
            }

            if (e.data === "update-categories") {
                onUpdateCategories();
            }
        }
    }
}