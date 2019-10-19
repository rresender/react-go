import React, { Component } from 'react';
import PropTypes from 'prop-types';
import UserList from './UserList.jsx';

class UserForm extends Component {
    onSubmit(e) {
        e.preventDefault();
        const node = this.refs.userName;
        const userName = node.value;
        this.props.setUserName(userName);
        node.value = '';
    }
    render() {
        return (
            <form onSubmit={this.onSubmit.bind(this)}>
                <div className='form-group'>
                    <input
                        ref='userName'
                        type='text'
                        className='form-control'
                        placeholder='Set Your Name...'
                    />
                </div>
            </form>
        );
    }
}

UserForm.propsType = {
    setUserName: PropTypes.func.isRequired
}

export default UserForm