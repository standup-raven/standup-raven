import React from 'react';

import logo from '../assets/images/logo.svg';

export class ChannelHeaderButtonIcon extends React.PureComponent {
    render() {
        return (
            <span
                dangerouslySetInnerHTML={{
                    __html: logo,
                }}
            />
        );
    }
}
