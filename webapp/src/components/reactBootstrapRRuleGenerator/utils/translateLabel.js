import {each, isFunction, isPlainObject, get} from 'lodash';

const replacePlaceholder = (text, replacements = {}) => {
    let result = text;
    each(replacements, (value, key) => {
        result = text.replace(`%{${key}}`, value);
    });

    return result;
};

const translateLabel = (translations, key, replacements = {}) => {
    if (isFunction(translations)) {
        return translations(key, replacements);
    } else if (isPlainObject(translations)) {
        return replacePlaceholder(
            get(translations, key, `[translation missing '${key}']`),
            replacements,
        );
    }

    return null;
};

export default translateLabel;
