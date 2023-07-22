// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function blurOnEscape(node: any) {
    const handleKey = (event: KeyboardEvent) => {
        if (event.key === 'Escape' && node) {
            if (typeof node.blur === 'function') node.blur();
            if (typeof node.value === 'string') node.value = '';
        }
    };

    node.addEventListener('keydown', handleKey);

    return {
        destroy() {
            node.removeEventListener('keydown', handleKey);
        }
    };
}
