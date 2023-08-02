import {MONTHS} from '../../../constants/index';

const computeYearlyOnMonth = (data, rruleObj) => {
    if (rruleObj.freq !== 0 || !rruleObj.bymonthday) {
        return data.repeat.yearly.on.month;
    }

    if (typeof rruleObj.bymonth === 'number') {
        return MONTHS[rruleObj.bymonth - 1];
    }

    return MONTHS[rruleObj.bymonth[0] - 1];
};

export default computeYearlyOnMonth;
