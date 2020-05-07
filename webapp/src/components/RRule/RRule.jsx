import * as React from 'react';
import PropTypes from 'prop-types';
import RRuleGenerator from 'react-rrule-generator';
import rruleStyles from 'react-rrule-generator/build/styles.css';
import {ControlLabel, FormGroup} from 'react-bootstrap';
import configModalStyles from '../configModal/style';
import style from './style.css';

import DatePicker from 'react-16-bootstrap-date-picker';

class RRule extends React.PureComponent {
    constructor(props) {
        super(props);
        this.state = RRule.getInitialState();

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

    rruleChangeHandler = (rrule) => {
        console.log(rrule);
        this.setState({
            rrule: rrule,
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
        this.props.onChange(rrule, startaDate);
    };

    render() {
        const style = configModalStyles.getStyle();

        return (
            <div>
                <FormGroup style={style.formGroup}>
                    <ControlLabel style={style.controlLabel}>
                        {'Start Date:'}
                    </ControlLabel>
                    {/*TODO add local formatted date in value*/}
                    <div className={'recurrence-start-date'}>
                        <DatePicker
                            id={'recurrence-start-date-picker'}
                            value={this.state.startDate}
                            onChange={this.startDateChangeHandler}
                        />
                    </div>
                </FormGroup>
                <FormGroup style={style.formGroup}>
                    <RRuleGenerator
                        onChange={this.rruleChangeHandler}
                        value={this.state.rrule}
                    />
                </FormGroup>
            </div>
        );
    }
}

RRule.PropTypes = {
    rrule: PropTypes.string,
    startDate: PropTypes.string,
    onChange: PropTypes.func.isRequired,
};

export default RRule;
