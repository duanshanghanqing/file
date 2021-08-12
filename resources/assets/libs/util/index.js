const $util = Object.create(null);

$util.sessionStorageGet = (key) => {
    let value = sessionStorage.getItem(key);
    try {
        value = JSON.parse(value)
    } catch (error) {
        // console.warn('sessionStorageGet', error)
    }
    return value;
}

$util.sessionStorageSet = (key, value) => {
    let _value = value;
    if (typeof value === 'object') {
        _value = JSON.stringify(value);
    }
    sessionStorage.setItem(key, _value);
}

$util.sessionStorageRemove = (key) => {
    sessionStorage.removeItem(key);
}

$util.sessionStorageClear = () => {
    sessionStorage.clear();
}

