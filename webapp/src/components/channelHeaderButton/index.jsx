import React from 'react';

import logo from '../../assets/images/logo.svg';
import './style.css';

class ChannelHeaderButtonIcon extends React.PureComponent {
    render() {
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

export default ChannelHeaderButtonIcon;
