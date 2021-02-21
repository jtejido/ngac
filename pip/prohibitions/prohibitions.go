package prohibitions

/**
 * Interface to maintain Prohibitions for an NGAC environment. This interface is in the common package because the
 * Prohibition service in the PDP will also implement this interface as well as any implementations in the PAP.
 */
type Prohibitions interface {
	/**
	 * Create a new prohibition.
	 */
	Add(*Prohibition)
	/**
	 * Get a list of all prohibitions
	 */
	All() []*Prohibition
	/**
	 * Retrieve a Prohibition and return the Object representing it.
	 */
	Get(string) *Prohibition
	/**
	 * Get all of the prohibitions a given entity is the direct subject of.  The subject can be a user, user attribute,
	 * or process.
	 */
	ProhibitionsFor(string) []*Prohibition
	/**
	 * Update the prohibition with the given name. Prohibition names cannot be updated.
	 */
	Update(string, *Prohibition)
	/**
	 * Delete the prohibition, and remove it from the data structure.
	 */
	Remove(string)
}
