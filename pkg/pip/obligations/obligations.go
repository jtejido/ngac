package obligations

type Obligations interface {
	/**
	 * Add the given obligation and enable it if the enable flag is true.
	 */
	Add(*Obligation, bool)
	/**
	 * Retrieves the obligation with the given label
	 */
	Get(string) *Obligation
	/**
	 * Return all obligations
	 */
	All() []*Obligation
	/**
	 * Update the obligation with the given label.  If the label in the provided object is not null and different from
	 * the label parameter, the label will also be updated.
	 */
	Update(string, *Obligation)
	/**
	 * Delete the obligation with the given label.
	 */
	Remove(string)
	/**
	 * Set the enable flag of the obligation with the given label.
	 */
	SetEnable(string, bool)
	/**
	 * Returns all enabled obligations
	 */
	GetEnabled() []*Obligation
}
