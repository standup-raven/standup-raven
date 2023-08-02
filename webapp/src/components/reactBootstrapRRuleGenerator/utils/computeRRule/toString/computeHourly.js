import {RRule} from 'rrule';

const computeHourly = ({interval}) => ({
    freq: RRule.HOURLY,
    interval,
});

export default computeHourly;
