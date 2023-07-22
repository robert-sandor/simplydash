import { dev } from '$app/environment';

const localDevHost = '192.168.1.90:8080';

function host(): string {
    return dev ? localDevHost : window.location.host;
}

export function baseUrl(): string {
    return `${window.location.protocol}//${host()}`;
}

export function websocketUrl(): string {
    const websocketProtocol = window.location.protocol === ':https' ? 'wss' : 'ws';
    return `${websocketProtocol}://${host()}/ws`;
}
