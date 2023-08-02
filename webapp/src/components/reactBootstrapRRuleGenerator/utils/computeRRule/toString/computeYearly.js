import {RRule} from 'rrule';

import computeYearlyOn from './computeYearlyOn';
import computeYearlyOnThe from './computeYearlyOnThe';

const computeYearly = ({mode, on, onThe}) => ({
    freq: RRule.YEARLY,
    ...(mode === 'on' ? computeYearlyOn(on) : computeYearlyOnThe(onThe)),
});

export default computeYearly;

