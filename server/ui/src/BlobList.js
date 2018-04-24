import React, { Component } from 'react';
import { connect } from 'react-redux';
import moment from 'moment';
import { copyLink } from './actions/index.js';

import ContentCopy from 'material-ui/svg-icons/content/content-copy';
import LinearProgress from 'material-ui/LinearProgress';
import Paper from 'material-ui/Paper';
import {List, ListItem} from 'material-ui/List';
import IconButton from 'material-ui/IconButton';

import {
  UPLOAD_START,
  UPLOAD_PROGRESS
} from './constants.js';

const expiresAt = (blob) => {
  return (
    <em>Expires {moment(blob.expires_at).from(moment().utc())}</em>
  );
};

const blobListItem = (blob) => {
  if (moment(blob.expires_at).isAfter(moment().utc())) {
    const buttons = (
      <div>
        <IconButton onClick={ copyLink(blob) }>
          <ContentCopy/>
        </IconButton>
      </div>
    );

    return (
      <ListItem
        key={blob.blob_id}
        disabled={true}
        primaryText={blob.filename}
        secondaryText={expiresAt(blob)}
        rightIconButton={buttons} />
    );
  }
};

const renderProgress = (upload) => {
  if (upload.status === UPLOAD_START) {
    return <LinearProgress />
  } else if (upload.status === UPLOAD_PROGRESS) {
    return <LinearProgress mode="determinate" min={0} max={100} value={upload.progress} />;
  }
}

class BlobList extends Component {
  renderUpload(upload) {
    return (
      <ListItem
        key={upload.id}
        disabled={true}
        primaryText={upload.file.name}
        secondaryText={renderProgress(upload)} />
    )
  }

  renderUploads() {
    const { uploads } = this.props;

    if (uploads) {
      return uploads.map(this.renderUpload);
    }
  }

  render() {
    const { blobs } = this.props;

    return (
      <Paper className="blobd-uploader-status">
        <List>
          {this.renderUploads()}
          {blobs.map(blobListItem)}
        </List>
      </Paper>
    );
  }
}

const mapStateToProps = (state, ownProps) => {
  console.log('state', state);

  return {
    uploads: state.upload,
    blobs: state.blobs,
  }
};

export default connect(mapStateToProps)(BlobList);
