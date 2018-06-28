import React, { Component } from 'react';
import { connect } from 'react-redux';
import moment from 'moment';
import { copyLink } from './actions/index.js';

import ContentCopy from 'material-ui/svg-icons/content/content-copy';
import LinearProgress from 'material-ui/LinearProgress';
import Paper from 'material-ui/Paper';
import {List, ListItem} from 'material-ui/List';
import IconButton from 'material-ui/IconButton';
import Divider from 'material-ui/Divider';

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

const uploadListItem = (upload) => {
  return (
    <ListItem
      key={upload.id}
      disabled={true}
      primaryText={upload.file.name}
      secondaryText={renderProgress(upload)} />
  );
}

const renderProgress = (upload) => {
  if (upload.status === UPLOAD_START) {
    return <LinearProgress />
  } else if (upload.status === UPLOAD_PROGRESS) {
    return (
      <LinearProgress
        mode="determinate"
        min={0}
        max={100}
        value={upload.progress} />
    );
  }
}

const maybeDivider = (uploadItems, blobItems) => {
  if (uploadItems.length && blobItems.length) {
    return <Divider />;
  }
}

class BlobList extends Component {
  renderList() {
    const { uploads, blobs } = this.props;

    const uploadItems = uploads.map(uploadListItem);
    const blobItems = blobs.map(blobListItem);

    return (
      <Paper className="blobd-uploader-status">
        <List>
          {uploadItems}
          {maybeDivider(uploadItems, blobItems)}
          {blobItems}
        </List>
      </Paper>
    )
  }

  render() {
    return this.renderList();
  }
}

const mapStateToProps = (state, ownProps) => {
  console.log('state', state);

  return {
    uploads: state.uploads,
    blobs: state.blobs,
  }
};

export default connect(mapStateToProps)(BlobList);
