import * as React from 'react';
import PropTypes from 'prop-types';
import {MenuItem, SplitButton} from 'react-bootstrap';
import style from './style.css';
import SentryBoundary from '../../SentryBoundary';

class TimePicker extends (SentryBoundary, React.PureComponent) {
    constructor(props) {
        super(props);
        this.state = TimePicker.getInitialState();

        // the sole purpose of this line is to prevent code formatters from marking the import as an unused one
        // TODO: see if we can configure this as a side effect instead of using it.
        // eslint-disable-next-line no-unused-vars
        const x = style;

        this.onChange = this.onChange.bind(this);
    }

    static get HOURS_MAX_VALUE() {
        // eslint-disable-next-line no-magic-numbers
        return 23;
    }

    static get MINUTES_MAX_VALUE() {
        // eslint-disable-next-line no-magic-numbers
        return 59;
    }

    static isPristine(time) {
        const pristineValues = [null, undefined, ''];

        return pristineValues.indexOf(time ? time.trim() : time) > -1;
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        if (nextProps.time === prevState.time && nextProps.bsStyle === prevState.bsStyle) {
            return null;
        }

        const time = nextProps.time || '00:00';
        const components = time.split(':');

        let hours;
        let minutes;

        if (components.length >= 2) {
            hours = components[0].trim();
            minutes = components[1].trim();
        }

        hours = parseInt(hours || 0, 10);
        hours = Math.max(0, Math.min(TimePicker.HOURS_MAX_VALUE, hours));
        // eslint-disable-next-line no-magic-numbers
        hours = hours < 10 ? '0' + hours : String(hours);

        minutes = parseInt(minutes || 0, 10);
        minutes = Math.max(0, Math.min(TimePicker.MINUTES_MAX_VALUE, minutes));
        // eslint-disable-next-line no-magic-numbers
        minutes = minutes < 10 ? '0' + minutes : String(minutes);

        return {
            time,
            hours,
            minutes,
            pristine: TimePicker.isPristine(),
            bsStyle: nextProps.bsStyle ? nextProps.bsStyle : 'default',
        };
    }

    static getInitialState() {
        return {
            time: '00:00',
            hours: '00',
            minutes: '00',
            pristine: true,
            bsStyle: 'default',
        };
    }

    onChange() {
        this.setState({
            pristine: false,
        });

        if (this.props.onChange) {
            const time = `${this.state.hours}:${this.state.minutes}`;
            this.props.onChange(time);
        }
    }

    render() {
        // eslint-disable-next-line no-shadow
        const style = getStyle();
        const hoursMenuItems = [];
        for (let i = 0; i <= TimePicker.HOURS_MAX_VALUE; ++i) {
            // eslint-disable-next-line no-magic-numbers
            const x = i < 10 ? '0' + i : String(i);
            hoursMenuItems.push(
                <MenuItem
                    key={i.toString()}
                    eventKey={x}
                    active={this.state.hours === x}
                >
                    {x}
                </MenuItem>);
        }

        const minutesMenuItems = [];
        for (let i = 0; i <= TimePicker.MINUTES_MAX_VALUE; ++i) {
            // eslint-disable-next-line no-magic-numbers
            const x = i < 10 ? '0' + i : String(i);
            minutesMenuItems.push(
                <MenuItem
                    key={i.toString()}
                    eventKey={x}
                    active={this.state.minutes === x}
                >
                    {x}
                </MenuItem>);
        }

        return (
            <div
                className={'time-picker'}
                style={{display: 'inline-block'}}
            >
                <SplitButton
                    id={`${this.props.id}-hours`}
                    bsStyle={this.state.bsStyle}
                    className={'hours'}
                    title={this.state.hours}
                    onSelect={(evt) => {
                        this.setState(
                            {hours: evt},
                            this.onChange,
                        );
                    }}
                >
                    {hoursMenuItems}
                </SplitButton>
                <span
                    className={'time-separator'}
                    style={style.timeSeparator}
                >{':'}</span>
                <SplitButton
                    id={`${this.props.id}-minutes`}
                    bsStyle={this.state.bsStyle}
                    className={'minutes'}
                    title={this.state.minutes}
                    onSelect={(evt) => {
                        this.setState(
                            {minutes: evt},
                            this.onChange,
                        );
                    }}
                >
                    {minutesMenuItems}
                </SplitButton>
            </div>
        );
    }
}

function getStyle() {
    return {
        timeSeparator: {
            paddingLeft: '5px',
            paddingRight: '5px',
        },
    };
}

TimePicker.propTypes = {
    id: PropTypes.string.isRequired,
    time: PropTypes.string,
    onChange: PropTypes.func,
    bsStyle: PropTypes.string,
};

export default TimePicker;
