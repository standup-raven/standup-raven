import React from 'react';

import logo from '../../assets/images/logo.svg';
import './style.css';
import PropTypes from 'prop-types';

class ChannelHeaderButtonIcon extends React.Component {
    getInitialState = () => {
        return {
            flag: true,
        };
    };

    constructor(props) {
        super(props);

        this.state = this.getInitialState();
    }

    render() {
        console.log('Renderring');
        console.log(this.props.channelID);
        console.log(this.state.flag);

        return (
            <span
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
};

export default ChannelHeaderButtonIcon;
