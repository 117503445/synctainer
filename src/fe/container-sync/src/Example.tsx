import { Combobox, ComboboxInput, ComboboxOption, ComboboxOptions } from '@headlessui/react'
import { useState } from 'react'
import { CheckIcon } from '@heroicons/react/20/solid'
type Person = {
    id: number;
    name: string;
};

const people = [
    { id: 1, name: 'Durward Reynolds' },
    { id: 2, name: 'Kenton Towne' },
    { id: 3, name: 'Therese Wunsch' },
    { id: 4, name: 'Benedict Kessler' },
    { id: 5, name: 'Katelyn Rohan' },
]

export default function Example() {
    const [selectedPerson, setSelectedPerson] = useState(people[0])
    const [query, setQuery] = useState('')

    const filteredPeople =
        query === ''
            ? people
            : people.filter((person) => {
                return person.name.toLowerCase().includes(query.toLowerCase())
            })

    const handleSelect = (value: Person | null) => {
        if (value !== null) {
            setSelectedPerson(value);
        }
    };

    return (
        <Combobox value={selectedPerson} onChange={handleSelect} onClose={() => setQuery('')}>
            <ComboboxInput
                aria-label="Assignee"
                displayValue={(person: Person) => person?.name}
                onChange={(event) => setQuery(event.target.value)}
            />
            <ComboboxOptions anchor="bottom" className="group flex gap-2 bg-white data-[focus]:bg-blue-100">
                {filteredPeople.map((person) => (
                    <ComboboxOption key={person.id} value={person} className="data-[focus]:bg-blue-100">
                        <CheckIcon className="invisible size-5 group-data-[selected]:visible" />
                        {person.name}
                    </ComboboxOption>
                ))}
            </ComboboxOptions>
        </Combobox>
    )
}