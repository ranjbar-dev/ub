self.addEventListener('install', e => {
  console.log('[custom service worker] installing', e);
});
self.addEventListener('activate', e => {
  console.log('[custom service worker] activating', e);
  return self.clients.claim();
});
self.addEventListener('fetch', e => {
  console.log('[custom service worker] fetching', e);
});
