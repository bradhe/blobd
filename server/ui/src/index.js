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

const shouldStoreState = (state) => {
  if (!state.uploads) {
    throw 'unexpected state object. expected uploads reducer.';
  }

  return state.uploads.length === 0;
};

const storeState = (state) => {
  // If there is no access to local storage then we can't cache the local
  // state. There's nothing to do!
  if (!window.localStorage) return;

  window.localStorage['$blobd.state'] = JSON.stringify(state);
};

const retrieveState = () => {
  // If there's no local storage available then there's nothing we can do here.
  if (!window.localStorage) return;

  const data = window.localStorage['$blobd.state'];

  // If there was nothing found then...well...
  if (!data) return {};

  try {
    return JSON.parse(data);
  } catch (e) {
    console.log('unable to retrieve state: ' + e);
  }

  return {};
};

const persistState = ({getState}) => {
  return (next) => (action) => {
    const ret = next(action);

    // We want to store the state as long as we are in a stable...state...
    const state = getState();

    if (shouldStoreState(state)) {
      storeState(state);
    }

    return ret;
  };
};

const store = createStore(
  reducer,
  retrieveState(),
  applyMiddleware(
    thunkMiddleware,
    persistState,
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
