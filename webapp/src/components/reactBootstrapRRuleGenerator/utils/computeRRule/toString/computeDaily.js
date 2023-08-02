import {RRule} from 'rrule';

const computeDaily = ({interval}) => ({
    freq: RRule.DAILY,
    interval,
});

export default computeDaily;
