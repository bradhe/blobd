import React, { Component } from 'react';
import { connect } from 'react-redux';
import './App.css';
import Uploader from './Uploader.js';
import BlobList from './BlobList.js';
import { MuiThemeProvider, getMuiTheme } from 'material-ui/styles';

import {
  fullBlack,
  white,
  grey800,
  deepPurple700,
  deepPurpleA700,
  deepPurpleA200,
  deepPurple200,
  deepPurple100
} from 'material-ui/styles/colors'

class App extends Component {
  constructor(props, context) {
    super(props, context);

    this.muiTheme = getMuiTheme({
      fontFamily: 'Montserrat',
      palette: {
	primary1Color: deepPurple700,
	primary2Color: deepPurple200,
	primary3Color: deepPurple100,
        accent1Color: deepPurpleA700,
        accent2Color: deepPurpleA200,
        accent3Color: deepPurpleA700,
	textColor: grey800,
	alternateTextColor: white,
	canvasColor: white,
	borderColor: deepPurple200,
	disabledColor: deepPurple100,
	pickerHeaderColor: deepPurpleA700,
	clockCircleColor: deepPurpleA200,
	shadowColor: fullBlack,
      },
      userAgent: props.userAgent
    });
  }

  render() {
    return (
      <MuiThemeProvider muiTheme={this.muiTheme}>
        <div className="App">
          <main className="blobd-container">
            <Uploader />
            <BlobList />
          </main>
        </div>
      </MuiThemeProvider>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  return {}
};

export default connect(mapStateToProps)(App);
