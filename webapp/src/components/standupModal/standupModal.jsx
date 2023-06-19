import React from 'react';
import PropTypes from 'prop-types';
import {Alert, Button, FormControl, FormGroup, InputGroup, Modal, OverlayTrigger, Tooltip} from 'react-bootstrap';
import request from 'superagent';
import Constants from '../../constants';
import reactStyles from './style';
import SentryBoundary from '../../SentryBoundary';
import * as HttpStatus from 'http-status-codes';
import Cookies from 'js-cookie';

const {formatText, messageHtmlToComponent} = window.PostUtils;

const standupModalCloseTimeout = 1000;
const standupTaskDefaultRowCount = 5;

class StandupModal extends (SentryBoundary, React.Component) {
    constructor(props) {
        super(props);
        this.state = StandupModal.getInitialState();
    }

    static get MODAL_CLOSE_TIMEOUT() {
        return standupModalCloseTimeout;
    }

    static get STANDUP_TASKS_DEFAULT_ROW_COUNT() {
        return standupTaskDefaultRowCount;
    }

    static getInitialState() {
        return {
            standup: {},
            activeTab: '',
            message: {
                show: false,
                text: '',
                type: 'info',
            },
            showSpinner: true,
            standupConfig: undefined,
            standupError: false,
            showStandupError: false,
            standupErrorMessage: '',
            standupErrorSubMessage: '',
        };
    }

    handleTasks = (key, e) => {
        const standup = {...this.state.standup};
        standup[key][e.target.name] = e.target.value;
        this.setState({
            standup,
        });
    };

    handleClose = () => {
        this.setState(StandupModal.getInitialState);
        this.props.close();
    };

    switchTabs = (direction) => {
        // this can be optimized by storing index rather than label in state variable
        const i = this.state.standupConfig.sections.indexOf(this.state.activeTab);
        const nextTab = this.state.standupConfig.sections[i + 1] || this.state.activeTab;
        const prevTab = this.state.standupConfig.sections[i - 1] || this.state.activeTab;

        this.setState({
            activeTab: direction === 'forward' ? nextTab : prevTab,
        });
    };

    handleSubmit = (event) => {
        event.preventDefault();

        const payload = this.prepareUserStandup();

        request
            .post(`${this.props.siteURL}/${Constants.URL_SUBMIT_USER_STANDUP}?channel_id=${this.props.channelID}`)
            .withCredentials()
            .send(payload)
            .set('X-CSRF-Token', Cookies.get(Constants.MATTERMOST_CSRF_COOKIE))
            .set('Content-Type', 'application/json')
            .end((err, res) => {
                if (err) {
                    this.setState({
                        message: {
                            show: true,
                            text: 'An error occurred while submitting standup.\n' + err.response.text,
                            type: 'danger',
                        },
                    });
                } else {
                    this.setState({
                        message: {
                            show: true,
                            text: 'Standup submitted successfully!',
                            type: 'success',
                        },
                    });
                    setTimeout(this.handleClose, StandupModal.MODAL_CLOSE_TIMEOUT);
                }
            });
    };

    getUserStandup = () => {
        return new Promise((resolve) => {
            request
                .get(`${this.props.siteURL}/${Constants.URL_SUBMIT_USER_STANDUP}?channel_id=${this.props.channelID}`)
                .withCredentials()
                .end((err, result) => {
                    if (result.ok) {
                        for (const sectionTitle of Object.keys(result.body.standup)) {
                            if (this.state.standup[sectionTitle] !== undefined) {
                                for (let i = 0; i < result.body.standup[sectionTitle].length; ++i) {
                                    this.state.standup[sectionTitle][`line${i + 1}`] = result.body.standup[sectionTitle][i];
                                }
                            }
                        }
                    } else if (result.status !== HttpStatus.NOT_FOUND) {
                        console.error(err);
                    }
                    resolve();
                });
            console.log('hello');
        });
    };

    getStandupConfig = () => {
        return new Promise((resolve) => {
            const url = `${this.props.siteURL}/${Constants.URL_STANDUP_CONFIG}?channel_id=${this.props.channelID}`;
            request
                .get(url)
                .withCredentials()
                .end((err, result) => {
                    if (result.ok) {
                        const state = {
                            standupConfig: result.body,
                            activeTab: result.body.sections[0],
                            standup: {},
                        };

                        result.body.sections.forEach((x) => {
                            state.standup[x] = {};
                        });
                        this.setState(state);
                    } else if (result.status !== HttpStatus.NOT_FOUND) {
                        console.error(err);
                    }
                    resolve();
                });
        });
    };

    prepareUserStandup = () => {
        const standup = {
            channelId: this.props.channelID,
            standup: {},
        };

        for (const sectionTitle of Object.keys(this.state.standup)) {
            standup.standup[sectionTitle] = Object.values(this.state.standup[sectionTitle])
                .map((x) => x.trim())
                .filter((x) => x !== '');
        }

        return standup;
    };

    componentDidUpdate(prevProp) {
        if (this.props.visible && !prevProp.visible) {
            this.getStandupConfig()
                .then(this.getUserStandup)
                .then(() => {
                    this.setState({showSpinner: false});
                })
                .catch(() => {
                    this.setState({showSpinner: false});
                });
        }
    }

    insertRows = (count, className, onChange) => {
        const rows = [];

        for (let i = 0; i <= Object.keys(this.state.standup[className] || {}).length; ++i) {
            rows.push(
                <FormGroup key={i.toString()}>
                    <InputGroup>
                        <InputGroup.Addon>{(i + 1) + '.'}</InputGroup.Addon>
                        <FormControl
                            type='text'
                            onChange={onChange}
                            name={'line' + (i + 1)}
                            className={className}
                            value={this.state.standup[className] ? (this.state.standup[className][`line${i + 1}`] || '') : ''}
                        />
                        <FormControl.Feedback/>
                    </InputGroup>
                </FormGroup>,
            );
        }

        return rows;
    };

    render() {
        const style = reactStyles.getStyle();

        let showStandupError = false;
        let standupErrorMessage = '';
        let standupErrorSubMessage = '';

        if (!this.state.standupConfig) {
            showStandupError = true;
            standupErrorMessage = 'Standup is not configured for this channel.';
            standupErrorSubMessage = 'Make sure you are filling the standup in the right channel or that standup has been configured in this channel.';
        } else if (!this.state.standupConfig.enabled) {
            showStandupError = true;
            standupErrorMessage = 'Standup is disabled for this channel.';
            standupErrorSubMessage = 'Please enable standup to continue using the features.';
        } else if (!this.state.standupConfig.members || this.state.standupConfig.members.length === 0) {
            showStandupError = true;
            standupErrorMessage = 'No members configured for this channel\'s standup.';
            standupErrorSubMessage = 'Please add some members to the standup to continue using the features.';
        } else if (this.state.standupConfig.members.indexOf(this.props.currentUserId) < 0) {
            showStandupError = true;
            standupErrorMessage = 'You are not a part of this channel\'s standup.';
            standupErrorSubMessage = 'Make sure you are filling standup in the right channel or that you were correctly added to the channel\'s standup.';
        } else if (this.props.isGuest) {
            showStandupError = true;
            standupErrorMessage = 'You are not allowed to submit standup.';
            standupErrorSubMessage = 'Guest users are not allowed to submit standup.';
        }

        const showSpinner = this.state.showSpinner;
        const showStandupForm = !showStandupError && this.state.standupConfig !== undefined;

        const sections = [];
        if (this.state.standupConfig && this.state.standupConfig.sections) {
            for (let i = 0; i < this.state.standupConfig.sections.length; ++i) {
                const sectionTitle = this.state.standupConfig.sections[i];
                sections.push(
                    <div
                        key={i.toString()}
                        id={sectionTitle}
                        className={this.state.activeTab === sectionTitle ? '' : 'hidden'}
                    >
                        {this.insertRows(StandupModal.STANDUP_TASKS_DEFAULT_ROW_COUNT, sectionTitle, (e) => {
                            this.handleTasks(sectionTitle, e);
                        })}
                    </div>,
                );
            }
        }

        const firstTab = this.state.standupConfig ? this.state.standupConfig.sections[0] : '';
        const lastTab = this.state.standupConfig ? this.state.standupConfig.sections[this.state.standupConfig.sections.length - 1] : '';

        return (
            <Modal
                show={this.props.visible}
                onHide={this.handleClose}
                backdrop={'static'}
            >

                <Modal.Header closeButton={true}>
                    <Modal.Title>
                        {Constants.PLUGIN_DISPLAY_NAME}
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body>
                    <div
                        className={showSpinner ? '' : 'hidden'}
                        style={style.spinner}
                    >
                        <img
                            src={`${this.props.siteURL}/${Constants.URL_SPINNER_ICON}`}
                            alt={'loading...'}
                        />
                    </div>

                    <span className={showSpinner ? 'hidden' : ''}>
                        <span className={showStandupError ? '' : 'hidden'}>
                            <span style={style.standupErrorMessage}>{standupErrorMessage}</span>
                            <br/><br/>
                            <span>{standupErrorSubMessage}</span>
                        </span>

                        <span className={showStandupForm ? '' : 'hidden'}>
                            <h5 style={style.header}>{messageHtmlToComponent(formatText(this.state.activeTab))}</h5>
                            <form style={style.form}>
                                <div className={'formContainer'}>
                                    {sections}
                                </div>
                            </form>
                            <Button
                                bsStyle='primary'
                                className={'fa fa-chevron-left'}
                                onClick={() => this.switchTabs('backward')}
                                disabled={this.state.activeTab === firstTab}
                                style={style.controlBtns}
                            />
                            <Button
                                bsStyle='primary'
                                className={'fa fa-chevron-right'}
                                onClick={() => this.switchTabs('forward')}
                                disabled={this.state.activeTab === lastTab}
                                style={style.controlBtns}
                            />
                        </span>
                    </span>
                </Modal.Body>
                <Modal.Footer>
                    <Button
                        type='button'
                        onClick={this.handleClose}
                        bsStyle='link'
                    >
                        {'Cancel'}
                    </Button>
                    <OverlayTrigger
                        placement={'bottom'}
                        overlay={
                            <Tooltip
                                id={'standup-submit-btn-tooltip'}
                                className={this.state.activeTab === lastTab ? 'hidden' : ''}
                            >
                                <strong>
                                    {'navigate to last tab to submit'}
                                </strong>
                            </Tooltip>
                        }
                    >
                        <Button
                            className={showStandupForm ? '' : 'hidden'}
                            type='submit'
                            bsStyle='primary'
                            onClick={this.handleSubmit}
                            disabled={this.state.activeTab !== lastTab}
                        >
                            {'Submit'}
                        </Button>
                    </OverlayTrigger>
                </Modal.Footer>
                <Alert
                    bsStyle={this.state.message.type}
                    style={style.alert}
                    className={(this.state.message.show ? '' : 'hidden')}
                >
                    {this.state.message.text}
                </Alert>

            </Modal>
        );
    }
}

StandupModal.propTypes = {
    channelID: PropTypes.string.isRequired,
    currentUserId: PropTypes.string.isRequired,
    close: PropTypes.func.isRequired,
    visible: PropTypes.bool.isRequired,
    siteURL: PropTypes.string.isRequired,
    isGuest: PropTypes.bool.isRequired,
};

export default StandupModal;
