namespace ASP.NETCoreWebAPI.Models
{
    public class ColonyCode
    {
        public int Id { get; set; } // Primary Key
        public int LobbyId { get; set; } // Attribute
        public int Value { get; set; } // Attribute - (Correct type?)
        public required Colony Colony { get; set; } // Foreign Key (Colony)
    }
}
