import type { VideoResult, VideoDetail, SearchResponse, DetailResponse } from '$lib/types';
import { browser } from '$app/environment';

const API_BASE = '/api';
const AUTH_HEADER = 'X-Auth-Password';
const STORAGE_KEY = 'searchav_auth_password';
const ADULT_MODE_KEY = 'searchav_adult_mode';

/** Get stored password */
function getPassword(): string | null {
	if (browser) {
		return localStorage.getItem(STORAGE_KEY);
	}
	return null;
}

/** Create headers with auth */
function createHeaders(): HeadersInit {
	const headers: HeadersInit = {};
	const password = getPassword();
	if (password) {
		headers[AUTH_HEADER] = password;
	}
	return headers;
}

/** Check if response is unauthorized */
function isUnauthorized(res: Response, data: { code: number }): boolean {
	return res.status === 401 || data.code === 401;
}

/** Search videos */
export async function search(query: string, includeAdult?: boolean): Promise<VideoResult[]> {
	const adult = includeAdult ?? getAdultMode();
	const params = new URLSearchParams({ q: query });
	if (adult) {
		params.set('adult', '1');
	}

	const res = await fetch(`${API_BASE}/search?${params.toString()}`, {
		headers: createHeaders()
	});

	const data: SearchResponse = await res.json();

	if (isUnauthorized(res, data)) {
		throw new AuthError('Unauthorized');
	}

	if (!res.ok || data.code !== 200) {
		throw new Error(data.msg || 'Search failed');
	}

	return data.list || [];
}

/** Get video detail */
export async function getDetail(sourceCode: string, vodId: number): Promise<VideoDetail> {
	const res = await fetch(`${API_BASE}/detail?source=${sourceCode}&id=${vodId}`, {
		headers: createHeaders()
	});

	const data: DetailResponse = await res.json();

	if (isUnauthorized(res, data)) {
		throw new AuthError('Unauthorized');
	}

	if (!res.ok || data.code !== 200) {
		throw new Error(data.msg || 'Failed to get detail');
	}

	return data.data;
}

/** Auth error class */
export class AuthError extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'AuthError';
	}
}

/** Save password to storage */
export function savePassword(password: string) {
	if (browser) {
		localStorage.setItem(STORAGE_KEY, password);
	}
}

/** Clear password from storage */
export function clearPassword() {
	if (browser) {
		localStorage.removeItem(STORAGE_KEY);
	}
}

/** Check if password exists */
export function hasPassword(): boolean {
	return !!getPassword();
}

/** Get adult mode setting */
export function getAdultMode(): boolean {
	if (browser) {
		return localStorage.getItem(ADULT_MODE_KEY) === '1';
	}
	return false;
}

/** Set adult mode setting */
export function setAdultMode(enabled: boolean) {
	if (browser) {
		if (enabled) {
			localStorage.setItem(ADULT_MODE_KEY, '1');
		} else {
			localStorage.removeItem(ADULT_MODE_KEY);
		}
	}
}
