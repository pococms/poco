if (document.readyState !== 'loading') {
    docReady();
} else {
    document.addEventListener('DOMContentLoaded', function () {
        docReady();
    });
}
