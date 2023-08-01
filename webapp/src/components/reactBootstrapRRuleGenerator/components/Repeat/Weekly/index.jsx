import React from 'react';
import PropTypes from 'prop-types';
import {toPairs} from 'lodash';
import numericalFieldHandler from '../../../utils/numericalFieldHandler';
import translateLabel from '../../../utils/translateLabel';

const RepeatWeekly = ({
    id,
    weekly: {
        interval,
        days,
        options,
    },
    handleChange,
    translations,
    weeklyFrequencyInputStyle,
}) => {
    let daysArray = toPairs(days);
    if (options.weekStartsOnSunday) {
        daysArray = daysArray.slice(-1).concat(daysArray.slice(0, -1));
    }

    return (
        <div>
            <div className='form-group d-flex align-items-sm-center'>
                <div className='col-sm-1 offset-sm-2 text-bold'>
                    {translateLabel(translations, 'repeat.weekly.every')}
                </div>
                <div>
                    <input
                        id={`${id}-interval`}
                        name='repeat.weekly.interval'
                        aria-label='Repeat weekly interval'
                        className='form-control'
                        value={interval}
                        onChange={numericalFieldHandler(handleChange)}
                        style={weeklyFrequencyInputStyle}
                    />
                </div>
                <div className='col-sm-1 d-flex align-items-center'>
                    {translateLabel(translations, 'repeat.weekly.weeks')}
                </div>
            </div>

            <div className='form-group row individual-day-of-week-selector'>
                <div className='btn-group btn-group-toggle offset-sm-2'>
                    {daysArray.map(([dayName, isDayActive]) => (
                        <label
                            htmlFor={`${id}-${dayName}`}
                            key={dayName}
                            className={`btn btn-primary day ${isDayActive ? 'active' : ''}`}
                        >
                            <input
                                type='checkbox'
                                id={`${id}-${dayName}`}
                                name={`repeat.weekly.days[${dayName}]`}
                                className='form-control'
                                checked={isDayActive}
                                onChange={(event) => {
                                    const editedEvent = {
                                        ...event,
                                        target: {
                                            ...event.target,
                                            value: !isDayActive,
                                            name: event.target.name,
                                        },
                                    };

                                    handleChange(editedEvent);
                                }}
                            />
                            {translateLabel(translations, `days_short.${dayName.toLowerCase()}`)}
                        </label>))
                    }
                </div>
            </div>
        </div>
    );
};

RepeatWeekly.propTypes = {
    id: PropTypes.string.isRequired,
    weekly: PropTypes.shape({
        interval: PropTypes.number.isRequired,
        days: PropTypes.shape({
            mon: PropTypes.bool.isRequired,
            tue: PropTypes.bool.isRequired,
            wed: PropTypes.bool.isRequired,
            thu: PropTypes.bool.isRequired,
            fri: PropTypes.bool.isRequired,
            sat: PropTypes.bool.isRequired,
            sun: PropTypes.bool.isRequired,
        }).isRequired,
        options: PropTypes.shape({
            weekStartsOnSunday: PropTypes.bool,
        }).isRequired,
    }).isRequired,
    handleChange: PropTypes.func.isRequired,
    translations: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
    weeklyFrequencyInputStyle: PropTypes.object,
};

export default RepeatWeekly;
