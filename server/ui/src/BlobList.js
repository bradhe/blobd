import React, { Component } from 'react';
import { connect } from 'react-redux';

import ContentCopy from 'material-ui/svg-icons/content/content-copy';
import FileDownload from 'material-ui/svg-icons/file/file-download';
import LinearProgress from 'material-ui/LinearProgress';
import Paper from 'material-ui/Paper';
import {List, ListItem} from 'material-ui/List';
import IconButton from 'material-ui/IconButton';

import {
  UPLOAD_START,
  UPLOAD_PROGRESS,
  UPLOAD_READY
} from './constants.js';

const blobListItem = (blob) => {
  const buttons = (
    <div>
      <IconButton>
        <ContentCopy/>
      </IconButton>
      <IconButton>
        <FileDownload/>
      </IconButton>
    </div>
  );

  return (
    <ListItem
      disabled={true}
      primaryText={blob.filename}
      rightIconButton={buttons} />
  );
};

const renderProgress = (upload) => {
  if (upload.status === UPLOAD_START) {
    return <LinearProgress />
  } else if (upload.status === UPLOAD_PROGRESS) {
    return <LinearProgress mode="determinate" min={0} max={100} value={upload.progress} />;
  }
}

class BlobList extends Component {
  renderCurrentUpload() {
    const { upload } = this.props;

    if (upload.status !== UPLOAD_READY) {
      return (
        <ListItem
          disabled={true}
          primaryText={upload.file.name}
          secondaryText={renderProgress(upload)} />
      )
    }
  }

  render() {
    const { blobs } = this.props;

    return (
      <Paper className="blobd-uploader-status">
        <List>
          {this.renderCurrentUpload()}
          {blobs.blobs.map(blobListItem)}
        </List>
      </Paper>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  return {
    upload: state.upload,
    blobs: state.blobs,
  }
};

export default connect(mapStateToProps)(BlobList);
