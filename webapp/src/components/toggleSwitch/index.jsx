import React from 'react';
import Switch from 'react-switch';
import SentryBoundary from '../../SentryBoundary';
import color from 'color';
import './style.css';

const darkenRatio = 0.3;

class ToggleSwitch extends (SentryBoundary, React.Component) {
    render() {
        return (
            <div className={'toggle-switch'} >
                <Switch
                    onChange={this.props.onChange}
                    checked={this.props.checked}
                    onColor={color(this.props.theme.linkColor).darken(darkenRatio).hex()}
                    boxShadow='0px 1px 5px rgba(0, 0, 0, 0.6)'
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
