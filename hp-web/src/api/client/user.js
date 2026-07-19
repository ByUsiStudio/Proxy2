import http from '../../data/http'

export function login(data) {
    return http({
        url: '/user/login',
        method: 'post',
        data
    })
}

export function register(data) {
    return http({
        url: '/user/register',
        method: 'post',
        data
    })
}

export function getSystemConfig() {
    return http({
        url: '/user/systemConfig',
        method: 'get'
    })
}

export function updateSystemConfig(data) {
    return http({
        url: '/user/updateSystemConfig',
        method: 'post',
        data
    })
}

export function updateUserStatus(params) {
    return http({
        url: '/user/updateUserStatus',
        method: 'get',
        params
    })
}
