import React, { Component } from 'react';
import { connect } from 'react-redux';
import './App.css';
import Uploader from './Uploader.js';
import BlobList from './BlobList.js';
import { MuiThemeProvider, getMuiTheme } from 'material-ui/styles';
import { fullBlack, white, purple500, grey800, deepPurpleA700, deepPurpleA200, deepPurple500, deepPurple300, deepPurple200, deepPurple100 } from 'material-ui/styles/colors'

class App extends Component {
  constructor(properties, context) {
    super(properties, context);

    this.muiTheme = getMuiTheme({
      palette: {
	primary1Color: deepPurple300,
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
      userAgent: properties.userAgent
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
