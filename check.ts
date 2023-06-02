interface Record {
	age: number,
	hobbies: hobby,
	id: number,
	name: string,
}

interface hobby {
	indoor: string[],
	outdoor: string[],
	wish: Wishes,
}

interface Wishes {
	current: string[],
}

