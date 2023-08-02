import React from 'react';
import PropTypes from 'prop-types';
import {MenuItem, SplitButton} from 'react-bootstrap';

import {DAYS} from '../../../constants/index';
import translateLabel from '../../../utils/translateLabel';

const RepeatMonthlyOnThe = ({
    id,
    mode,
    onThe,
    hasMoreModes,
    handleChange,
    translations,
    monthlyOnTheDayDropdownStyle,
}) => {
    const isActive = mode === 'on the';

    return (
        <div className={`form-group d-flex align-items-sm-center text-bold ${!isActive && 'opacity-50'}`}>
            <div className={'combo-label'}>
                <div className='col-sm-1'>
                    {hasMoreModes && (
                        <input
                            id={id}
                            type='radio'
                            name='repeat.monthly.mode'
                            aria-label='Repeat monthly on the'
                            value='on the'
                            checked={isActive}
                            onChange={handleChange}
                        />
                    )}
                </div>
                <div className='col-sm-8'>
                    {translateLabel(translations, 'repeat.monthly.on_the')}
                </div>
            </div>

            <div className='week-of-month-selector'>
                <SplitButton
                    bsStyle={'link'}
                    id={`${id}-which`}
                    name='repeat.monthly.onThe.which'
                    aria-label='Repeat monthly on the which'
                    key={onThe.which}
                    title={onThe.which}
                    disabled={!isActive}
                    style={monthlyOnTheDayDropdownStyle}
                    dropup={true}
                    onSelect={(eventKey) => handleChange({target: {name: 'repeat.monthly.onThe.which', value: eventKey}})}
                >
                    <MenuItem eventKey='First'>{translateLabel(translations, 'numerals.first')}</MenuItem>
                    <MenuItem eventKey='Second'>{translateLabel(translations, 'numerals.second')}</MenuItem>
                    <MenuItem eventKey='Third'>{translateLabel(translations, 'numerals.third')}</MenuItem>
                    <MenuItem eventKey='Fourth'>{translateLabel(translations, 'numerals.fourth')}</MenuItem>
                    <MenuItem eventKey='Last'>{translateLabel(translations, 'numerals.last')}</MenuItem>
                </SplitButton>

            </div>

            <div className='day-of-week-selector'>
                <SplitButton
                    bsStyle={'link'}
                    id={`${id}-day`}
                    name='repeat.monthly.onThe.day'
                    aria-label='Repeat monthly on the day'
                    key={onThe.day}
                    title={onThe.day}
                    disabled={!isActive}
                    style={monthlyOnTheDayDropdownStyle}
                    dropup={true}
                    onSelect={(eventKey) => handleChange({target: {name: 'repeat.monthly.onThe.day', value: eventKey}})}
                >
                    {DAYS.map((day) => (
                        <MenuItem
                            eventKey={day}
                            key={day}
                        >
                            {translateLabel(translations, `days.${day.toLowerCase()}`)}
                        </MenuItem>
                    ))}
                </SplitButton>
            </div>

        </div>
    );
};

RepeatMonthlyOnThe.propTypes = {
    id: PropTypes.string.isRequired,
    mode: PropTypes.oneOf(['on', 'on the']).isRequired,
    onThe: PropTypes.shape({
        which: PropTypes.oneOf(['First', 'Second', 'Third', 'Fourth', 'Last']).isRequired,
        day: PropTypes.oneOf(DAYS).isRequired,
    }).isRequired,
    hasMoreModes: PropTypes.bool.isRequired,
    handleChange: PropTypes.func.isRequired,
    translations: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
    monthlyOnTheDayDropdownStyle: PropTypes.object,
};

export default RepeatMonthlyOnThe;
