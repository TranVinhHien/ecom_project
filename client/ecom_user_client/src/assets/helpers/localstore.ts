function setDataToLocalStore(key: string, duLieu: any) {
    return new Promise((resolve, reject) => {
        try {
            localStorage.setItem(key, JSON.stringify(duLieu));
            resolve("localStorage.setItem.success");
        } catch (error) {
            reject(error);
        }
    });
}
function getDataToLocalStore<T>(key: string): T | undefined {
    const duLieu = localStorage.getItem(key);
    if (!duLieu) {
        return undefined;
    }
    try {
        return JSON.parse(duLieu) as T;
    } catch (error) {
        return undefined;
    }
}
export { setDataToLocalStore, getDataToLocalStore };

