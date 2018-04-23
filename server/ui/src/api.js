import request from 'superagent';
import Promise from 'bluebird';

const DEFAULT_HOST = 'http://localhost:8765'

const apiRoot = (uri) => {
  let url = DEFAULT_HOST;

  if (uri) {
    url += uri;
  }

  return url;
};

function responseHandler(reject, resolve, err, res, skipNotFound = false) {
  const error = err || res.error;

  if (res && res.text) {
    res.body = JSON.parse(res.text);
  }

  const errorMsg = error && res && res.body ? res.body : error;

  return errorMsg ? reject(errorMsg) : resolve(res.body);
}


const api = {
  send() {
    return new Promise((resolve, reject) => {
      request.post(apiRoot('/api/turn-off'))
        .end((err, res) => {
          responseHandler(reject, resolve, err, res);
        });
    });
  },

  getBlob(id, jwt) {
    return new Promise((resolve, reject) => {
      request.get(apiRoot(`/${id}`))
        .setHeader('Authorization', `Bearer ${jwt}`)
        .end((err, res) => {
          responseHandler(reject, resolve, err, res);
        });
    });
  },

  upload(file) {
    return new Promise((resolve, reject) => {
      console.log('uploading', file);
      request.post(apiRoot('/'))
        .attach('content', file, {mimeType: file.type})
        .end((err, res) => {
          responseHandler(reject, resolve, err, res);
        });
    });
  }
}

export default api;
