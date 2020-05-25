import * as React from 'react';
import PropTypes from 'prop-types';
import RRuleGenerator from 'react-rrule-generator';
import rruleStyles from 'react-rrule-generator/build/styles.css';
import {ControlLabel, FormGroup} from 'react-bootstrap';
import configModalStyles from '../configModal/style';
import style from './style.css';
import reactStyles from './style';

import DatePicker from 'react-16-bootstrap-date-picker';

class RRule extends React.PureComponent {
    constructor(props) {
        super(props);
        this.state = RRule.getInitialState();
        this.configModalReactStyles = configModalStyles.getStyle();

        // eslint-disable-next-line no-unused-vars
        const x = rruleStyles;

        // eslint-disable-next-line no-unused-vars
        const y = style;
    }

    static getInitialState = () => {
        return {
            rrule: 'RRULE:FREQ=WEEKLY;INTERVAL=1;BYDAY=MO,TU,WE,TH,FR',
            startDate: new Date().toISOString(),
            startDateDisplay: new Date(),
        };
    };

    static getDerivedStateFromProps(nextProps, prevState) {
        // dont update state if nothing changed
        if (nextProps.rrule === prevState.rrule || nextProps.startDate === prevState.startDate) {
            return null;
        }

        return {
            rrule: nextProps.rrule || prevState.rrule,
            startDate: nextProps.startDate || prevState.startDate,
        };
    }

    componentDidMount() {
        this.sendChanges(this.state.rrule, this.state.startDate);
    }

    rruleChangeHandler = (rrule) => {
        console.log(rrule);
        this.setState({
            rrule,
        });

        this.sendChanges(rrule, this.state.startDate);
    };

    startDateChangeHandler = (isoDate) => {
        this.setState({
            startDate: isoDate,
        });

        this.sendChanges(this.state.rrule, isoDate);
    };

    sendChanges = (rrule, startaDate) => {
        this.props.onChange(rrule.replace('RRULE:', ''), startaDate);
    };

    render() {
        return (
            <div>
                <FormGroup style={this.configModalReactStyles.formGroup}>
                    <ControlLabel style={this.configModalReactStyles.controlLabel}>
                        {'Start Date:'}
                    </ControlLabel>
                    {/*TODO add local formatted date in value*/}
                    <div
                        className={'recurrence-start-date'}
                        style={reactStyles.getStyle().recurrenceDatepicker}
                    >
                        <DatePicker
                            id={'recurrence-start-date-picker'}
                            value={this.state.startDate}
                            onChange={this.startDateChangeHandler}
                            showClearButton={false}
                            style={{width: '60%'}}
                        />
                    </div>
                </FormGroup>
                <FormGroup
                    style={this.configModalReactStyles.formGroup}
                    className={'standup-recurrence'}
                >
                    <RRuleGenerator
                        config={{
                            hideStart: true,
                            hideEnd: true,
                        }}
                        onChange={this.rruleChangeHandler}
                        value={this.state.rrule}
                        customCalendar={DatePicker}
                    />
                </FormGroup>
            </div>
        );
    }
}

RRule.propTypes = {
    rrule: PropTypes.string,
    startDate: PropTypes.string,
    onChange: PropTypes.func.isRequired,
};

export default RRule;
