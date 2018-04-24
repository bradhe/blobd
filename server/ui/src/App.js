import React, { Component } from 'react';
import { connect } from 'react-redux';
import './App.css';
import Uploader from './Uploader.js';
import { MuiThemeProvider, getMuiTheme } from 'material-ui/styles';
import { fullBlack, red700, red800, red50, red100, red500, pinkA200, darkBlack, white, grey100, grey500, grey300 } from 'material-ui/styles/colors'

class App extends Component {
  constructor(properties, context) {
    super(properties, context);

    this.muiTheme = getMuiTheme({
      palette: {
	primary1Color: red500,
	primary2Color: red100,
	primary3Color: red50,
        accent1Color: pinkA200,
        accent2Color: grey100,
        accent3Color: grey500,
	textColor: darkBlack,
	alternateTextColor: white,
	canvasColor: white,
	borderColor: grey300,
	disabledColor: red100,
	pickerHeaderColor: red500,
	clockCircleColor: darkBlack,
	shadowColor: fullBlack,
      },
      userAgent: properties.userAgent
    });
  }

  render() {
    return (
      <MuiThemeProvider muiTheme={this.muiTheme}>
        <div className="App">
          <main className="container">
            <Uploader {...this.props.upload} />
          </main>
        </div>
      </MuiThemeProvider>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  return {
    upload: state.upload,
  }
};

export default connect(mapStateToProps)(App);
