import React from 'react';
import PropTypes from 'prop-types';
import {ControlLabel, MenuItem, SplitButton} from 'react-bootstrap';

import translateLabel from '../../utils/translateLabel';

import RepeatYearly from './Yearly/index';
import RepeatMonthly from './Monthly/index';
import RepeatWeekly from './Weekly/index';
import RepeatDaily from './Daily/index';
import RepeatHourly from './Hourly/index';

const Repeat = ({
    id,
    repeat: {
        frequency,
        yearly,
        monthly,
        weekly,
        daily,
        hourly,
        options,
    },
    handleChange,
    translations,
    repeatDropdownStyle,
    monthlyFrequencyInputStyle,
    weeklyFrequencyInputStyle,
    monthlyOnDayDropdownStyle,
    monthlyOnTheDayDropdownStyle,
}) => {
    const isOptionAvailable = (option) => !options.frequency || options.frequency.indexOf(option) !== -1;
    const isOptionSelected = (option) => frequency === option;

    return (
        <div>
            <div className='form-group repeat-frequency'>
                <ControlLabel>
                    {translateLabel(translations, 'repeat.label')}
                </ControlLabel>
                <SplitButton
                    bsStyle={'link'}
                    title={frequency}
                    key={frequency}
                    id={'dropdown-basic'}
                    onSelect={(eventKey) => handleChange({target: {name: 'repeat.frequency', value: eventKey}})}
                    name={'repeat.frequency'}
                    style={repeatDropdownStyle}
                >
                    {isOptionAvailable('Yearly') &&
                    <MenuItem eventKey='Yearly'>{translateLabel(translations, 'repeat.yearly.label')}</MenuItem>}
                    {isOptionAvailable('Monthly') &&
                    <MenuItem eventKey='Monthly'>{translateLabel(translations, 'repeat.monthly.label')}</MenuItem>}
                    {isOptionAvailable('Weekly') &&
                    <MenuItem eventKey='Weekly'>{translateLabel(translations, 'repeat.weekly.label')}</MenuItem>}
                    {isOptionAvailable('Daily') &&
                    <MenuItem eventKey='Daily'>{translateLabel(translations, 'repeat.daily.label')}</MenuItem>}
                    {isOptionAvailable('Hourly') &&
                    <MenuItem eventKey='Hourly'>{translateLabel(translations, 'repeat.hourly.label')}</MenuItem>}
                </SplitButton>
            </div>
            <div className={'repeat-configuration'}>
                {
                    isOptionSelected('Yearly') &&
                    <RepeatYearly
                        id={`${id}-yearly`}
                        yearly={yearly}
                        handleChange={handleChange}
                        translations={translations}
                    />
                }
                {
                    isOptionSelected('Monthly') &&
                    <RepeatMonthly
                        id={`${id}-monthly`}
                        monthly={monthly}
                        handleChange={handleChange}
                        translations={translations}
                        monthlyFrequencyInputStyle={monthlyFrequencyInputStyle}
                        monthlyOnDayDropdownStyle={monthlyOnDayDropdownStyle}
                        monthlyOnTheDayDropdownStyle={monthlyOnTheDayDropdownStyle}
                    />
                }
                {
                    isOptionSelected('Weekly') &&
                    <RepeatWeekly
                        id={`${id}-weekly`}
                        weekly={weekly}
                        handleChange={handleChange}
                        translations={translations}
                        weeklyFrequencyInputStyle={weeklyFrequencyInputStyle}
                    />
                }
                {
                    isOptionSelected('Daily') &&
                    <RepeatDaily
                        id={`${id}-daily`}
                        daily={daily}
                        handleChange={handleChange}
                        translations={translations}
                    />
                }
                {
                    isOptionSelected('Hourly') &&
                    <RepeatHourly
                        id={`${id}-hourly`}
                        hourly={hourly}
                        handleChange={handleChange}
                        translations={translations}
                    />
                }
            </div>
        </div>
    );
};

Repeat.propTypes = {
    id: PropTypes.string.isRequired,
    repeat: PropTypes.shape({
        frequency: PropTypes.oneOf(['Yearly', 'Monthly', 'Weekly', 'Daily', 'Hourly']).isRequired,
        yearly: PropTypes.object.isRequired,
        monthly: PropTypes.object.isRequired,
        weekly: PropTypes.object.isRequired,
        daily: PropTypes.object.isRequired,
        hourly: PropTypes.object.isRequired,
        options: PropTypes.shape({
            frequency: PropTypes.arrayOf(PropTypes.oneOf(['Yearly', 'Monthly', 'Weekly', 'Daily', 'Hourly'])),
            yearly: PropTypes.oneOf(['on', 'on the']),
            monthly: PropTypes.oneOf(['on', 'on the']),
        }).isRequired,
    }).isRequired,
    handleChange: PropTypes.func.isRequired,
    translations: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
    repeatDropdownStyle: PropTypes.object,
    weeklyFrequencyInputStyle: PropTypes.object,
    monthlyFrequencyInputStyle: PropTypes.object,
    monthlyOnDayDropdownStyle: PropTypes.object,
    monthlyOnTheDayDropdownStyle: PropTypes.object,
};

export default Repeat;
