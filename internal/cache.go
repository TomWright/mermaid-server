package internal

// DiagramCache provides the ability to cache diagram results.
type DiagramCache interface {
	// Store stores a diagram in the cache.
	Store(diagram *Diagram) error
	// Has returns true if we have a cache stored for the given diagram description.
	Has(diagram *Diagram) (bool, error)
	// Get returns a cached version of the given diagram description.
	Get(diagram *Diagram) (*Diagram, error)
}

// NewDiagramCache returns an implementation of DiagramCache.
func NewDiagramCache() DiagramCache {
	return &inMemoryDiagramCache{
		idToDiagram: map[string]*Diagram{},
	}
}

// inMemoryDiagramCache is an in-memory implementation of DiagramCache.
type inMemoryDiagramCache struct {
	idToDiagram map[string]*Diagram
}

// Store stores a diagram in the cache.
func (c *inMemoryDiagramCache) Store(diagram *Diagram) error {
	id, err := diagram.ID()
	if err != nil {
		return err
	}
	c.idToDiagram[id] = diagram
	return nil
}

// Has returns true if we have a cache stored for the given diagram description.
func (c *inMemoryDiagramCache) Has(diagram *Diagram) (bool, error) {
	id, err := diagram.ID()
	if err != nil {
		return false, err
	}
	if d, ok := c.idToDiagram[id]; ok && d != nil {
		return true, nil
	}
	return false, nil
}

// Get returns a cached version of the given diagram description.
func (c *inMemoryDiagramCache) Get(diagram *Diagram) (*Diagram, error) {
	id, err := diagram.ID()
	if err != nil {
		return nil, err
	}
	if d, ok := c.idToDiagram[id]; ok && d != nil {
		return d, nil
	}
	return nil, nil
}
