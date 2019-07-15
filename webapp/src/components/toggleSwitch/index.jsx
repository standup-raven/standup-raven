import React from 'react';
import Switch from 'react-switch';
import SentryBoundary from '../../SentryBoundary';
import color from 'color';
import './style.css';

class ToggleSwitch extends (SentryBoundary, React.Component) {
    render() {
        return (
            <div className={'toggle-switch'} >
                <Switch
                    onChange={this.props.onChange}
                    checked={this.props.checked}
                    onColor={color(this.props.theme.linkColor).darken(0.3).hex()}
                    offColor={color(this.props.theme.centerChannelColor).lighten(0.9).hex()}
                    offHandleColor={color(this.props.theme.centerChannelColor).lighten(0.7).hex()}
                    onHandleColor={this.props.theme.linkColor}
                    handleDiameter={23}
                    uncheckedIcon={false}
                    checkedIcon={false}
                    height={13}
                    width={35}
                />
            </div>
        );
    }
}

export default ToggleSwitch;
