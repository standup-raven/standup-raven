import React from 'react';

import logo from '../../assets/images/logo.svg';
import './style.css';
import PropTypes from 'prop-types';
import RavenClient from '../../raven-client';

class ChannelHeaderButtonIcon extends React.Component {
    constructor(props) {
        super(props);

        this.myRef = React.createRef();
        this.state = this.getInitialState();
    }

    componentDidMount() {
        RavenClient.Config.getActiveChannels(this.props.siteURL)
            .then((activeChannels) => {
                const activeChannelMap = {};
                activeChannels.forEach((x) => {
                    activeChannelMap[x] = true;
                });
                this.setState({
                    activeChannels: activeChannelMap,
                });
            });
    }

    componentDidUpdate(prevProps, prevState, snapshot) {
        if (prevProps.added !== this.props.added || prevProps.removed !== this.props.removed) {
            const activeChannels = this.state.activeChannels;

            if (this.props.added !== prevProps.added) {
                // new active channel is added
                activeChannels[this.props.added] = true;
            }

            if (this.props.removed !== prevProps.removed) {
                // new channel was removed
                activeChannels[this.props.removed] = undefined;
            }

            this.setState({
                activeChannels,
            });
        }
    }

    getInitialState = () => {
        return {
            activeChannels: {},
            parent: undefined,
        };
    };

    handleRef = (ref) => {
        if (ref) {
            this.setState({
                parent: ref.parentNode,
            });
        }
    }

    isChannelHeaderButtonInDropdown = () => {
        const classList = this.state.parent.parentNode.parentNode.parentNode.parentNode.classList;
        return classList.contains('dropdown') && classList.contains('btn-group');
    }

    getIconParentToHide = () => {
        if (this.isChannelHeaderButtonInDropdown()) {
            return this.state.parent.parentNode.parentNode;
        }
        return this.state.parent;
    }

    render() {
        if (this.state.parent) {
            const targetParent = this.getIconParentToHide();
            if (this.state.activeChannels[this.props.channelID]) {
                targetParent.classList.remove('hidden');
            } else {
                targetParent.classList.add('hidden');
            }
        }

        return (
            <span
                ref={this.handleRef}
                className={'raven-icon'}
                dangerouslySetInnerHTML={{
                    __html: logo,
                }}
            />
        );
    }
}

ChannelHeaderButtonIcon.propTypes = {
    channelID: PropTypes.string.isRequired,
    siteURL: PropTypes.string.isRequired,
    added: PropTypes.string.isRequired,
    removed: PropTypes.string.isRequired,
};

export default ChannelHeaderButtonIcon;
