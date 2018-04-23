export const upload = (file) => ({
  type: 'UPLOAD',
  brightness: file,
});

export const download = (id, token) => ({
  type: 'DOWNLOAD',
  id: id,
  token: token
});
