import React, { Component } from 'react';
import { connect } from 'react-redux';
import { upload } from './actions';

class FileUpload extends Component {
  constructor(props, context) {
    super(props, context);

    this.onFileChange = this.onFileChange.bind(this);
  }

  onFileChange(e) {
    const { dispatch } = this.props;

    // TODO: Should there be support for multiple files?
    let file = e.target.files[0];
    dispatch(upload(file));
  }

  render() {
    return (
      <div className="blobd-file-uploade">
        <input type="file" onChange={this.onFileChange} />
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  return {}
};

export default connect(mapStateToProps)(FileUpload);
