import React, { Component } from 'react';
import { connect } from 'react-redux';
import { upload } from './actions';
import { PropTypes } from 'prop-types';

import FloatingActionButton from 'material-ui/FloatingActionButton';
import ContentAdd from 'material-ui/svg-icons/content/add';

const propTypes = {
  status: PropTypes.string,
  file: PropTypes.object,
  progress: PropTypes.number,
};

class Uploader extends Component {
  constructor(props, context) {
    super(props, context);

    this.onFileChange = this.onFileChange.bind(this);
    this.onUploadButtonClick = this.onUploadButtonClick.bind(this);

    this.fileInput = React.createRef();
  }

  onUploadButtonClick(e) {
    this.fileInput.current.click();
  }

  onFileChange(e) {
    const { dispatch } = this.props;

    if (e.target.files) {
      const files = [...e.target.files];
      files.forEach((file) => dispatch(upload(file)));
    }
  }

  render() {
    return (
      <div className="blobd-uploader">
        <div className="blobd-uploader-button">
          <input ref={this.fileInput} type="file" id="blobd-uploader-file" style={{display: 'none'}} onChange={this.onFileChange}/>

          <FloatingActionButton onClick={this.onUploadButtonClick}>
            <ContentAdd />
          </FloatingActionButton>
        </div>
      </div>
    );
  }
}

Uploader.propTypes = propTypes;

const mapStateToProps = (state, ownProps) => {
  return {}
};

export default connect(mapStateToProps)(Uploader);
