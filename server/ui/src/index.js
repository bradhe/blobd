import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import { createStore, applyMiddleware } from 'redux'
import { Provider } from 'react-redux'
import thunkMiddleware from 'redux-thunk'
import reducer from './reducers'
import './index.css';
import injectTapEventPlugin from 'react-tap-event-plugin';

const store = createStore(
  reducer,
  applyMiddleware(
    thunkMiddleware
  )
);

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
);

registerServiceWorker();
injectTapEventPlugin();
