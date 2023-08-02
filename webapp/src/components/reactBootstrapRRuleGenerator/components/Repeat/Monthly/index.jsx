import React from 'react';
import PropTypes from 'prop-types';

import numericalFieldHandler from '../../../utils/numericalFieldHandler';
import translateLabel from '../../../utils/translateLabel';

import RepeatMonthlyOn from './On';
import RepeatMonthlyOnThe from './OnThe';

const RepeatMonthly = ({
    id,
    monthly: {
        mode,
        interval,
        on,
        onThe,
        options,
    },
    handleChange,
    translations,
    monthlyFrequencyInputStyle,
    monthlyOnDayDropdownStyle,
    monthlyOnTheDayDropdownStyle,
}) => {
    const isTheOnlyOneMode = (option) => options.modes === option;
    const isOptionAvailable = (option) => !options.modes || isTheOnlyOneMode(option);

    return (
        <div>
            <div className='form-group d-flex align-items-sm-center'>
                <div className='col-sm-1 offset-sm-2 text-bold'>
                    {translateLabel(translations, 'repeat.monthly.every')}
                </div>
                <div>
                    <input
                        id={`${id}-interval`}
                        name='repeat.monthly.interval'
                        aria-label='Repeat monthly interval'
                        className='form-control'
                        value={interval}
                        onChange={numericalFieldHandler(handleChange)}
                        style={monthlyFrequencyInputStyle}
                    />
                </div>
                <div className='col-sm-1 d-flex align-items-center'>
                    {translateLabel(translations, 'repeat.monthly.months')}
                </div>
            </div>

            {isOptionAvailable('on') && (
                <RepeatMonthlyOn
                    id={`${id}-on`}
                    mode={mode}
                    on={on}
                    hasMoreModes={!isTheOnlyOneMode('on')}
                    handleChange={handleChange}
                    translations={translations}
                    monthlyOnDayDropdownStyle={monthlyOnDayDropdownStyle}

                />
            )}
            {isOptionAvailable('on the') && (
                <RepeatMonthlyOnThe
                    id={`${id}-onThe`}
                    mode={mode}
                    onThe={onThe}
                    hasMoreModes={!isTheOnlyOneMode('on the')}
                    handleChange={handleChange}
                    translations={translations}
                    monthlyOnTheDayDropdownStyle={monthlyOnTheDayDropdownStyle}
                />
            )}

        </div>
    );
};

RepeatMonthly.propTypes = {
    id: PropTypes.string.isRequired,
    monthly: PropTypes.shape({
        mode: PropTypes.oneOf(['on', 'on the']).isRequired,
        interval: PropTypes.number.isRequired,
        on: PropTypes.object.isRequired,
        onThe: PropTypes.object.isRequired,
        options: PropTypes.shape({
            modes: PropTypes.oneOf(['on', 'on the']),
        }).isRequired,
    }).isRequired,
    handleChange: PropTypes.func.isRequired,
    translations: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
    monthlyFrequencyInputStyle: PropTypes.object,
    monthlyOnDayDropdownStyle: PropTypes.object,
    monthlyOnTheDayDropdownStyle: PropTypes.object,
};

export default RepeatMonthly;
