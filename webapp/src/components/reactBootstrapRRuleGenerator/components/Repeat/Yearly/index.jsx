// This file is unused. We are saving this for future use. Reference: https://github.com/standup-raven/react-bootstrap-rrule-generator

import React from 'react';
import PropTypes from 'prop-types';

import RepeatYearlyOn from './On';
import RepeatYearlyOnThe from './OnThe';

const RepeatYearly = ({
    id,
    yearly: {
        mode,
        on,
        onThe,
        options,
    },
    handleChange,
    translations,
}) => {
    const isTheOnlyOneMode = (option) => options.modes === option;
    const isOptionAvailable = (option) => !options.modes || isTheOnlyOneMode(option);
    return (
        <div>
            {isOptionAvailable('on') && (
                <RepeatYearlyOn
                    id={`${id}-on`}
                    mode={mode}
                    on={on}
                    hasMoreModes={!isTheOnlyOneMode('on')}
                    handleChange={handleChange}
                    translations={translations}
                />
            )}
            {isOptionAvailable('on the') && (
                <RepeatYearlyOnThe
                    id={`${id}-onThe`}
                    mode={mode}
                    onThe={onThe}
                    hasMoreModes={!isTheOnlyOneMode('on the')}
                    handleChange={handleChange}
                    translations={translations}
                />
            )}
        </div>
    );
};

RepeatYearly.propTypes = {
    id: PropTypes.string.isRequired,
    yearly: PropTypes.shape({
        mode: PropTypes.oneOf(['on', 'on the']).isRequired,
        on: PropTypes.object.isRequired,
        onThe: PropTypes.object.isRequired,
        options: PropTypes.shape({
            modes: PropTypes.oneOf(['on', 'on the']),
        }).isRequired,
    }).isRequired,
    handleChange: PropTypes.func.isRequired,
    translations: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
};

export default RepeatYearly;
