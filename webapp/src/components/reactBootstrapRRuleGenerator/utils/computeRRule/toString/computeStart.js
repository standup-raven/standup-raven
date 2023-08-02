import moment from 'moment';

const computeStart = ({onDate: {date}}) => {
    let result;
    // verify that incoming date is valid
    // by seeing if it can be converted into a moment object.
    // if not, then create a new date
    if (!moment.isMoment(moment(date))) {
        result = new Date().setMilliseconds(0);
    }

    return {
        dtstart: moment(result).toDate(),
    };
};

export default computeStart;
