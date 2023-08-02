import {RRule} from 'rrule';

const computeOptions = ({hideStart, weekStartsOnSunday}) => {
    const options = {};

    if (hideStart) {
        options.dtstart = null;
    }

    if (weekStartsOnSunday) {
        options.wkst = RRule.SU;
    }

    return options;
};

export default computeOptions;
