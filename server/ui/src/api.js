import request from 'superagent';
import Promise from 'bluebird';

const DEFAULT_HOST = 'http://localhost:8765';

const getHost = () => {
  return window.BLOBD_HOST || DEFAULT_HOST;
}

const apiRoot = (uri) => {
  let url = getHost();

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


export function download(id, jwt) {
  return new Promise((resolve, reject) => {
    request.get(apiRoot(`/${id}`))
      .setHeader('Authorization', `Bearer ${jwt}`)
      .end((err, res) => {
        responseHandler(reject, resolve, err, res);
      });
  });
};

export function uploadWithProgress(file, progress) {
  return new Promise((resolve, reject) => {
    request.post(apiRoot('/'))
      .attach('content', file, {mimeType: file.type})
      .on('progress', (event) => {
        progress(event.percent);
      })
      .end((err, res) => {
        responseHandler(reject, resolve, err, res);
      });
  });
};
