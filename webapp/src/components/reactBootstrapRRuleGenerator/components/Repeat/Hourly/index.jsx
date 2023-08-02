// This file is unused. We are saving this for future use. Reference: https://github.com/standup-raven/react-bootstrap-rrule-generator

import React from 'react';
import PropTypes from 'prop-types';

import numericalFieldHandler from '../../../utils/numericalFieldHandler';
import translateLabel from '../../../utils/translateLabel';

const RepeatHourly = ({
    id,
    hourly: {
        interval,
    },
    handleChange,
    translations,
}) => (
    <div className='form-group row d-flex align-items-sm-center'>
        <div className='col-sm-1 offset-sm-2'>
            {translateLabel(translations, 'repeat.hourly.every')}
        </div>
        <div className='col-sm-2'>
            <input
                id={`${id}-interval`}
                name='repeat.hourly.interval'
                aria-label='Repeat hourly interval'
                className='form-control'
                value={interval}
                onChange={numericalFieldHandler(handleChange)}
            />
        </div>
        <div className='col-sm-1'>
            {translateLabel(translations, 'repeat.hourly.hours')}
        </div>
    </div>
);

RepeatHourly.propTypes = {
    id: PropTypes.string.isRequired,
    hourly: PropTypes.shape({
        interval: PropTypes.number.isRequired,
    }).isRequired,
    handleChange: PropTypes.func.isRequired,
    translations: PropTypes.oneOfType([PropTypes.object, PropTypes.func]).isRequired,
};

export default RepeatHourly;
