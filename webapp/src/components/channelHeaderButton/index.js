import {connect} from 'react-redux';

import util from '../../utils';
import Selectors from '../../selectors';

import ChannelHeaderButtonIcon from './button';

const mapStateToProps = (state) => ({
    channelID: state.entities.channels.currentChannelId,
    siteURL: util.trimTrailingSlash(util.getValueSafely(state, 'entities.general.config.SiteURL', '')),
    added: Selectors.addedActiveChannel(state),
    removed: Selectors.removedActiveChannel(state),
});

export default connect(mapStateToProps)(ChannelHeaderButtonIcon);
